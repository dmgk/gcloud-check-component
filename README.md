Check for a new Google Cloud SDK component version and generate component manifest.

#### Installation

    go get github.com/dmgk/gcloud-check-component

#### Usage

Check for a new version:

    $ gcloud-check-component app-engine-go-linux-x86_64 20181102165100
	Found new build for app-engine-go-linux-x86_64: 20181102165140 (have 20181102165100)

Generate manifest:

    $ gcloud-check-component -manifest app-engine-go
    {
      "components": [
        {
          "dependencies": [
            "app-engine-go-darwin-x86",
            "app-engine-go-darwin-x86_64",
            "app-engine-go-linux-x86",
            "app-engine-go-linux-x86_64",
            "app-engine-go-windows-x86",
            "app-engine-go-windows-x86_64",
            "app-engine-python",
            "core"
          ],
          "details": {
            "description": "Provides the tools to develop Go applications on App Engine.",
            "display_name": "App Engine Go Extensions"
          },
          "id": "app-engine-go",
          "is_configuration": false,
          "is_hidden": false,
          "is_required": false,
          "platform": {},
          "version": {
            "build_number": 0,
            "version_string": ""
          }
        }
      ],
      "revision": 20190215163830,
      "schema_version": {
        "no_update": false,
        "url": "https://dl.google.com/dl/cloudsdk/channels/rapid/google-cloud-sdk.tar.gz",
        "version": 3
      },
      "version": "235.0.0"
    }
