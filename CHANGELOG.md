# Changelog
All notable changes to this project will be documented in this file.

## 2.1.0 - 2023-04-05
* Added update campaigns and feature flag configurations instantaneously with Real-Time Streaming Architecture: [`Documentation`](https://developers.kameleoon.com/go-sdk.html#streaming) or [`Product Updates`](https://www.kameleoon.com/en/blog/real-time-streaming)
* Added a new method [`OnUpdateConfiguration`](https://developers.kameleoon.com/go-sdk.html#onupdateconfiguration) to handle events when configuration data is updated in real time
* Added a new method [`GetEngineTrackingCode`](https://developers.kameleoon.com/go-sdk.html#GetEngineTrackingCode) which can be used to simplify utilization of hybrid mode
* Minor bug fixes

## 2.0.6 - 2023-03-24
* Minor bug fixing for `is among the values` operator for [`CustomData`](https://developers.kameleoon.com/go-sdk.html#customdata).
* Renaming of methods:
    - `RetrieveDataFromRemoteSource`-> [`GetRemoteData`](https://developers.kameleoon.com/go-sdk.html#GetRemoteData)

## 2.0.5 - 2023-03-13
* Added possibility for [`CustomData`](https://developers.kameleoon.com/go-sdk.html#customdata) to use variable argument list of values
* Fixed issue with [`TriggerExperiment`](https://developers.kameleoon.com/go-sdk.html#triggerexperiment) returning wrong error for client with no experiments
* Fixed issue with provided `Config` with not specified values; related to [initialization](https://developers.kameleoon.com/go-sdk.html#1-initialize-the-kameleoon-client)

## 2.0.4 - 2023-03-07
* [`GetVariationAssociatedData`](https://developers.kameleoon.com/go-sdk.html#getvariationassociateddata) Fixed. No need to unquote bytes upon obtaining anymore.

## 2.0.3 - 2023-02-27
* Minor bug fixing

## 2.0.2 - 2023-01-02
* Removed dependency on first version

## 2.0.1 - 2023-01-02
* Fixed issue with distribution of v2

## 2.0.0 - 2023-01-02
* Significantly improved configuration load time
* Added support for **Experiment** & **Exclusive Campaign** conditions. Related to [`TriggerExperiment`](https://developers.kameleoon.com/go-sdk.html#triggerexperiment)
* Renaming of methods:
    - `ActivateFeature`-> [`IsFeatureActive`](https://developers.kameleoon.com/go-sdk.html#IsFeatureActive)
    - `ErrNotActivated` -> `ErrNotAllocated`. Related to [`TriggerExperiment`](https://developers.kameleoon.com/go-sdk.html#triggerexperiment)
* Methods added for obtaining experiment and feature flag lists along with feature variables:
    - [`GetFeatureAllVariables`](https://developers.kameleoon.com/go-sdk.html#GetFeatureAllVariables)
    - [`GetFeatureList`](https://developers.kameleoon.com/go-sdk.html#GetFeatureList)
    - [`GetActiveFeatureListForVisitor`](https://developers.kameleoon.com/go-sdk.html#GetActiveFeatureListForVisitor)
    - [`GetExperimentList`](https://developers.kameleoon.com/go-sdk.html#GetExperimentList)
    - [`GetExperimentListForVisitor`](https://developers.kameleoon.com/go-sdk.html#GetExperimentListForVisitor)
* Added support of `is among the values` operator for Custom Data
* Added KameleoonData [`Device`](https://developers.kameleoon.com/go-sdk.html#device) data. Possible values are: **Phone**, **Tablet**, **Desktop**.
* Removed KameleoonData `Interest`

## 1.0.6 - 2022-04-12
* Added method for retrieving data from remote source: [`RetrieveDataFromRemoteSource`](https://developers.kameleoon.com/go-sdk.html#retrievedatafromremotesource)

## 1.0.5 - 2022-02-28
* Added support of multi-environment for feature flags, Related to [`ActivateFeature`](https://developers.kameleoon.com/go-sdk.html#activatefeature), [`GetFeatureVariable`](https://developers.kameleoon.com/go-sdk.html#getfeaturevariable)
* Added checking for status of site_code (Enable / Disable). Related to [`ActivateFeature`](https://developers.kameleoon.com/go-sdk.html#activatefeature), [`TriggerExperiment`](https://developers.kameleoon.com/go-sdk.html#triggerexperiment)

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
