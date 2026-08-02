apiVersion: v1
type: cli/v1
name: "cm-delete"
version: "@VERSION@"
runtime: subprocess
sourceURL: https://github.com/runzhliu/helm-delete
config:
  usage: cm-delete
  shortHelp: "Delete a specific version of a chart from a ChartMuseum repository"
runtimeConfig:
  platformCommand:
    - command: ${HELM_PLUGIN_DIR}/bin/helm-cm-delete
    - os: windows
      command: ${HELM_PLUGIN_DIR}\bin\helm-cm-delete.exe
  platformHooks:
    install:
      - command: sh
        args:
          - ${HELM_PLUGIN_DIR}/scripts/install_plugin.sh
    update:
      - command: sh
        args:
          - ${HELM_PLUGIN_DIR}/scripts/install_plugin.sh
