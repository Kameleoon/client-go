# Changelog
All notable changes to this project will be documented in this file.

## 3.0.2 - 2023-11-27
### Bug fixes
* Stability and performance improvements

## 3.0.1 - 2023-11-24
### Bug fixes
* Stability and performance improvements

## 3.0.0 - 2023-11-24
### Breaking changes
* Increased the minimum required version of the Go language to [1.18](https://go.dev/doc/go1.18).
* Renamed `Client` to `KameleoonClient`
* Removed `NewClient` function. Instead, use [`KameleoonClientFactory`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#kameleoonclientfactory)
* Removed `RunWhenReady` method. Instead, use [`WaitInit`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#waitinit)
* Changed `Config`:
    - Renamed `Config` to `KameleoonClientConfig`
    - Changed the default request timeout to `10` seconds
    - Renamed the `Timeout` and `timeout` configuration fields to `DefaultTimeout` and [`default_timeout`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk#additional-configuration), respectively)
    - Renamed the `ConfigUpdateInterval` and `config_update_interval` configuration fields to `RefreshInterval` and [`refresh_interval`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk#additional-configuration), respectively.)
    - `LoadConfig` function and `KameleoonClientConfig.Load` can now return an error.
* The `Cfg` field is no longer accessible.
* Renamed `GetFeatureAllVariables` method to `GetFeatureVariationVariables`(https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#getfeaturevariationvariables)
* Removed all methods and errors related to **experiments**:
  * Methods:
    - `TriggerExperiment`
    - `GetVariationAssociatedData`
    - `GetExperimentList`
    - `GetExperimentListForVisitor`
  * Error types:
    - `ErrExperimentConfigNotFound`
    - `ErrNotTargeted`
    - `ErrNotAllocated`
    - `ErrSiteCodeDisabled`
* Changed errors:
  * Moved error types into `errs` package
  * Added `UnexpectedStatusCode` error, which can be thrown by the following methods:
    - [`GetRemoteData`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#getremotedata)
    - [`GetRemoteVisitorData`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#getremotevisitordata)
  * Renamed the following errors:
    - `ErrFeatureConfigNotFound` to `FeatureNotFound`
    - `ErrFeatureVariableNotFound` to `FeatureVariableNotFound`
    - `ErrVariationNotFound` to `FeatureVariationNotFound`
    - `ErrCredentialsNotFound` to `ConfigCredentialsInvalid`
    - `ErrVisitorCodeNotValid` to `VisitorCodeInvalid`
* The new error `FeatureEnvironmentDisabled` may be returned by the following methods:
    - [`GetFeatureVariationKey`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#getfeaturevariationkey)
    - [`GetFeatureVariable`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#getfeaturevariable)
    - [`GetFeatureVariationVariables`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#getfeaturevariationvariables)
* Changed `Data` types:
  * [`Browser`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#browser):
    - Added constructor [`NewBrowser`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#newbrowser)
    - Hid the `Type` and `Version` fields and replaced them with getter methods
  * [`Conversion`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#conversion):
    - Added constructors [`NewConversion`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#newconversion) and [`NewConversionWithRevenue`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#newconversionwithrevenue)
    - Hid the `GoalId`, `Revenue`, and `Negative` fields and replaced them with getter methods
  * [`CustomData`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#customdata):
    - Changed the data type of the `ID` field to `int`
    - Hid the `ID` field and replaced with a getter method
    - Renamed `GetValues` to `Values`
  * [`Device`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#device):
    - Added constructor [`NewDevice`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#newdevice)
    - Hid the `Type` field and replaced with a getter method
  * [`PageView`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#pageview):
    - Added constructors [`NewPageView`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#newpageview), [`NewPageViewWithTitle`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#newpageviewwithtitle)
    - Hid the `URL`, `Title`, and `Referrers` fields and replaced with the getters
    - Changed `url` from optional to a required parameter
  * [`UserAgent`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#useragent):
    - Added constructor [`NewUserAgent`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#newuseragent)
    - Hid `Value` field and replaced with the getter
* Reworked cookies:
    - Removed `SetVisitorCode` method
    - Removed `ObtainVisitorCode` method
    - Removed parameter `topLevelDomain` from [`GetVisitorCode`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#getvisitorcode). Instead, use the `top_level_domain` parameter in the [configuration](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#additional-configuration)
    - Made [`GetVisitorCode`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#getvisitorcode) return `(string, error)` pair.
* Removed methods that were deprecated in 2.x versions:
    - `RetrieveDataFromRemoteSource`
* Removed visitor data max size:
    - Removed `visitor_data_max_size` configuration field
    - Removed `VisitorDataMaxSize` field from `Config`

### Features
* Added a [`SetLegalConsent`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#setlegalconsent) method to determine the types of data Kameleoon includes in tracking requests. This helps you adhere to legal and regulatory requirements while responsibly managing visitor data. You can find more information in the [Consent management policy](https://help.kameleoon.com/consent-management-policy/).
* Implemented `KameleoonClientFactory`(https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#kameleoonclientfactory) to manage `KameleoonClient` instances:
    - [`Create`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/create)
    - [`CreateFromFile`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#createfromfile)
    - [`Forget`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#forget)
* Added [`WaitInit`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#waitinit) method to wait until the initialization process is completed
* Added new parameters for `KameleoonClientConfig`:
    - `SessionDuration` ([`session_duration`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk#additional-configuration) configuration field accordingly)
    - `TopLevelDomain` field to `Config` ([`top_level_domain`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk#additional-configuration) configuration field accordingly)

### Bug fixes
* Stability and performance improvements

## 2.3.1 - 2023-10-03
### Features
* Added support for older versions of the Go language. You can now use version 1.12 or later (previously, the minimum version was 1.16).

## 2.3.0 - 2023-09-12
### Features
* Added [`GetRemoteVisitorData`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#GetRemoteVisitorData) method to fetch a visitor's remote data (with an optional capability to add the fetched data to the visitor).
### Bug fixes
* Stability and performance improvements

## 2.2.0 - 2023-08-20
* Stability and performance improvements

## 2.1.1 - 2023-06-28
* Added new conditions for targeting:
    - `Visitor Code`
    - `SDK Language`
    - [`Page Title & Page Url`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#pageview)
    - [`Browser`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#browser)
    - [`Device`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#device)
    - [`Conversion`](https://developers.kameleoon.com/feature-management-and-experimentation/web-sdks/go-sdk/#trackconversion)

## 2.1.0 - 2023-04-05
* Added update campaigns and feature flag configurations instantaneously with Real-Time Streaming Architecture: [`Documentation`](https://developers.kameleoon.com/go-sdk.html#streaming) or [`Product Updates`](https://www.kameleoon.com/en/blog/real-time-streaming)
* Added a new methods:
    - [`OnUpdateConfiguration`](https://developers.kameleoon.com/go-sdk.html#onupdateconfiguration) to handle events when configuration data is updated in real time
    - [`GetEngineTrackingCode`](https://developers.kameleoon.com/go-sdk.html#GetEngineTrackingCode) which can be used to simplify utilization of hybrid mode
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
