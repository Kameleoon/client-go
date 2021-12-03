#!/bin/sh
# deploy - push selected files to github

echo "<------- START CHECKING ENVIRONMENT ------->"
if [ -z "${GITHUB_TOKEN}" ]; then
    echo "Missing GITHUB_TOKEN environment variable"
    echo "<------- FAILED CHECKING ENVIRONMENT ------->"
    exit 1
fi

deploy_directory="client-go"
github_account_name="Kameleoon"
email_sdk="sdk@kameleoon.com"
github_repo_url="https://${GITHUB_TOKEN}@github.com/${github_account_name}/${deploy_directory}.git"
echo "<------- SUCCESS CHECKING ENVIRONMENT ------->"

echo "Prepare for deployment"
# Init github repository inside deploy folder
rm -rf "${deploy_directory}"

# Clone git and get version
git config --global user.email "${email_sdk}"
git clone "${github_repo_url}"
cd ${deploy_directory}
version=$(echo $(git describe --tags) | awk -F. -v OFS=. '{$NF++;print}')
while [ "$1" != "" ]; do
    case $1 in
    -v | --version)
        version=$2
        ;;
    esac
    shift
done
cd ../
rm -rf "${deploy_directory}"/*

# Copy needed files for deploy into a directory
tar cf deploy.tar --exclude="${deploy_directory}" *
mv deploy.tar "${deploy_directory}"/deploy.tar
cd "${deploy_directory}"
tar xf deploy.tar
rm deploy.tar
find . -name "*test*" -type f -delete
rm -rf *test*

echo "Deploying version ${version}"

# Commit and push new files
git add *
git commit -m "GO SDK ${version}"
git push --force

# Create tag and push
git tag "v${version}" main
git push origin "v${version}"

# Remove deploy folder
cd ../
rm -rf "${deploy_directory}"

echo "Finished to deploy version ${version}"
