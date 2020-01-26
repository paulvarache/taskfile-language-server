curl -sL https://taskfile.dev/install.sh | sh
p="$(pwd)/bin"
echo "##vso[task.prependpath]$p"
