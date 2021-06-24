#!/bin/sh
# deploy - push selected files to github

deploy_directory="client-go"
github_account_name="Kameleoon"


echo "Prepare for deployment"
# Init github repository inside deploy folder
rm -rf "${deploy_directory}"

# Clone git and get version
git config --global user.email "ggrandjean@kameleoon.com"
git clone "git@github.com:${github_account_name}/${deploy_directory}.git"
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
git commit -m "Version ${version}"
git push --force

# Create tag and push
git tag "${version}" main
git push origin "${version}"

# Remove deploy folder
cd ../
rm -rf "${deploy_directory}"

echo "Finished to deploy version ${version}"
