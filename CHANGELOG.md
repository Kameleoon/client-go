# Changelog
All notable changes to this project will be documented in this file.

## 1.0.4 - 2021-12-10
* Added scheduling functionality for [`ActivateFeature`](https://developers.kameleoon.com/go-sdk.html#activatefeature)
* Added VisitorCodeNotValid exception when empty or exceeding the limit of 255 chars. Related to [`ActivateFeature`](https://developers.kameleoon.com/go-sdk.html#activatefeature), [`TriggerExperiment`](https://developers.kameleoon.com/go-sdk.html#triggerexperiment), [`AddData`](https://developers.kameleoon.com/go-sdk.html#adddata), [`TrackConversion`](https://developers.kameleoon.com/go-sdk.html#trackconversion), [`FlushVisitor`](https://developers.kameleoon.com/go-sdk.html#flush)

## 1.0.3 - 2021-12-06
* GraphQL API is using now instead of REST
* Improved SDK stability

## 1.0.2 - 2021-12-03
* Fixed issue with wrong bucketing. Related to [`TriggerExperiment`](https://developers.kameleoon.com/go-sdk.html#triggerexperiment)

## 1.0.1 - 2021-11-30
* Improved SDK stability

## 1.0.0 - 2021-06-24
* Added Fasthttp
