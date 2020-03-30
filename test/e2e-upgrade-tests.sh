#!/usr/bin/env bash

# Copyright 2019 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script runs the end-to-end tests against Knative Eventing
# Operator built from source.  It is started by prow for each PR. For
# convenience, it can also be executed manually.

# If you already have a Knative cluster setup and kubectl pointing
# to it, call this script with the --run-tests arguments and it will use
# the cluster and run the tests.

# Calling this script without arguments will create a new cluster in
# project $PROJECT_ID, start knative in it, run the tests and delete the
# cluster.

source $(dirname $0)/../vendor/knative.dev/test-infra/scripts/e2e-tests.sh

# Latest eventing operator release.
readonly LATEST_EVENTING_OPERATOR_RELEASE_VERSION=$(git tag | sort -V | tail -1)
readonly LATEST_EVENTING_RELEASE_VERSION="v0.13.3"

OPERATOR_DIR=$(dirname $0)/..
KNATIVE_EVENTING_DIR=${OPERATOR_DIR}/..

# Namespace used for tests
readonly TEST_NAMESPACE="operator-tests"

declare -A COMPONENTS
COMPONENTS=(
  ["eventing.yaml"]="config"
  ["in-memory-channel.yaml"]="config/channels/in-memory-channel"
)

function install_eventing_operator() {
  header "Installing Knative Eventing operator"

  # Deploy the operator
  ko apply -f config/
  wait_until_pods_running default || fail_test "Eventing Operator did not come up"
}

function install_latest_operator_release() {
  header "Installing Knative Eventing operator latest public release"
  local full_url="https://github.com/knative/eventing-operator/releases/download/${LATEST_EVENTING_OPERATOR_RELEASE_VERSION}/eventing-operator.yaml"

  local release_yaml="$(mktemp)"
  wget "${full_url}" -O "${release_yaml}" \
      || fail_test "Unable to download latest Knative Eventing Operator release."

  kubectl apply -f "${release_yaml}" || fail_test "Knative Eventing Operator latest release installation failed"
  create_custom_resource
  wait_until_pods_running ${TEST_NAMESPACE}
}

function create_custom_resource() {
  echo ">> Creating the custom resource of Knative Eventing:"
  cat <<EOF | kubectl apply -f -
apiVersion: operator.knative.dev/v1alpha1
kind: KnativeEventing
metadata:
  name: knative-eventing
  namespace: ${TEST_NAMESPACE}
EOF
}

function knative_setup() {
  echo ">> Creating test namespaces"
  kubectl create namespace $TEST_NAMESPACE
  install_latest_operator_release
}

function install_head() {
  install_eventing_operator
}

function generate_latest_eventing_manifest() {
  # Go the directory to download the source code of knative eventing
  cd ${KNATIVE_EVENTING_DIR}

  # Download the source code of knative eventing
  git clone https://github.com/knative/eventing.git
  cd eventing
  COMMIT_ID=$(git rev-parse --verify HEAD)
  echo ">> The latest commit ID of Knative Eventing is ${COMMIT_ID}."
  mkdir -p output

  # Generate the manifest
  LABEL_YAML_CMD=(cat)
  local all_yamls=()
  for yaml in "${!COMPONENTS[@]}"; do
    local config="${COMPONENTS[${yaml}]}"
    echo "Building Knative Eventing - ${config}"
    ko resolve -P -f ${config}/ | "${LABEL_YAML_CMD[@]}" > ${yaml}
  done

  EVENTING_YAML=${KNATIVE_EVENTING_DIR}/eventing/eventing.yaml
  IMC_YAML=${KNATIVE_EVENTING_DIR}/eventing/in-memory-channel.yaml
  if [[ -f "${EVENTING_YAML}" && -f "${IMC_YAML}" ]]; then
    echo ">> Replacing the current manifest in operator with the generated manifest"
    rm -rf ${OPERATOR_DIR}/cmd/manager/kodata/knative-eventing/*
    cp ${EVENTING_YAML} ${OPERATOR_DIR}/cmd/manager/kodata/knative-eventing/eventing.yaml
    cp ${IMC_YAML} ${OPERATOR_DIR}/cmd/manager/kodata/knative-eventing/eventing-imc.yaml
  else
    echo ">> The eventing.yaml and in-memory-channel.yaml were not generated, so keep the current manifest"
  fi

  # Go back to the directory of operator
  cd ${OPERATOR_DIR}
}

# Script entry point.

# Skip installing istio as an add-on
initialize $@ --skip-istio-addon

TIMEOUT=20m

install_head

# If we got this far, the operator installed Knative Eventing
header "Running tests for Knative Eventing Operator"
failed=0

# Run the postupgrade tests
go_test_e2e -tags=postupgrade -timeout=${TIMEOUT} ./test/upgrade || failed=1

# Verify with the bash script to make sure there is no resource with the label of the previous release.
list_resources="deployment,pod,service,cm,crd,sa,ClusterRole,ClusterRoleBinding,ValidatingWebhookConfiguration,\
MutatingWebhookConfiguration,Secret,RoleBinding"
result="$(kubectl get ${list_resources} -l eventing.knative.dev/release=${LATEST_EVENTING_RELEASE_VERSION} --all-namespaces 2>/dev/null)"

# If the ${result} is not empty, we fail the tests, because the resources from the previous release still exist.
if [[ ! -z ${result} ]] ; then
  header "The following obsolete resources still exist:"
  echo "${result}"
  fail_test "The resources with the label of previous release have not been removed."
fi

# Require that tests succeeded.
(( failed )) && fail_test

success
