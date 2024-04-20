# Copyright 2024 Ross Light
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# SPDX-License-Identifier: Apache-2.0

{
  description = "Helper for creating a Vault plugin directory";

  outputs = { ... }: {
    lib.mkPluginDirectory = { plugins, pkgs }:
      let
        inherit (builtins) map;
        inherit (pkgs.lib.meta) getExe;
        inherit (pkgs.lib.strings) concatLines escapeShellArg getName optionalString;

        registerScriptName = "register-vault-plugins.sh";

        pluginScript = { binary
          , type ? "secret"
          , pname ? getName binary
          , version ? binary.version or ""
          }:
          let
            exePath = getExe binary;
            command = pname + (optionalString (version != "") "-${version}");
          in ''
            ln -s ${escapeShellArg exePath} "$out/bin/${command}"
            sum="$(sha256sum ${escapeShellArg exePath} | cut --delimiter=' ' --fields=1)"
            echo vault plugin register -sha256="$sum" -command=${escapeShellArg command} -version=${escapeShellArg version} ${escapeShellArg type} ${escapeShellArg pname} >> "$out/bin/${registerScriptName}"
          '';
      in
        pkgs.runCommandLocal "vault-plugin-directory" {} (''
          mkdir -p "$out/bin"
          echo ${escapeShellArg ("#!" + pkgs.runtimeShell)} > "$out/bin/${registerScriptName}"
          echo 'set -euo pipefail' >> "$out/bin/${registerScriptName}"
          ${concatLines (map pluginScript plugins)}
          chmod +x "$out/bin/${registerScriptName}"
        '');
  };
}
