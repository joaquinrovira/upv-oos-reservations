<center>

# UPV Office of Sports Reservations Agent

<img src="./.img/UPV.jpg"  height="100" />
<img src="./.img/muscles-clipart-ghoper.gif" height="100"  />

</center>

An attempt at automating the weekly reservations for the Universitat Politècnica de València's Office of Sports, written in Go.

## Why

In July of 2022 the Office of Sports announces that the university's sports facilities will be freely available for all students and associated personnel. Although all activities are free there is a limited availability for each. Some activities follow a first-come first-served policy. Others however, follow a weekly open call for reservations, beggining every Saturday at 10:00. For highly demanded activities, they run out of available slots quickly (in a matter of minutes) and people may cancel at any time during the week (up to 1 hour before the activity starts).

In order to reduce the time spent checking for reservation, I decided to build this tool. It automates the reservation process. Also works as an excuse to learn a new programming language I've had my eye on for a while. 

## Building the agent

Install `go` following the [official instructions](https://go.dev/doc/install) and build the binary with 

```bash
# Clone the repo
git clone https://github.com/joaquinrovira/upv-oos-reservations
cd upv-oos-reservations

# Build the binary
go build .

# Run the executable
./upv-oos-reservations
``` 

## How to configure the agent

There are two main aspects of the application to configure. Firstly, the environment variables. They can either be set normally **or by writing to a `.env` file**. Environment variable take precedence over `.env` values. For a full set of configurable environment variables, check out [`vars.go`](./lib/vars/vars.go).

<center>

| Variable            | Description                                                                                                  |
| ------------------- | ------------------------------------------------------------------------------------------------------------ |
| `UPV_USER`          | Intranet username                                                                                            |
| `UPV_PASS`          | Intranet password                                                                                            |
| `UPV_ACTIVITY_TYPE` | Internal activity type (see section below for more info)                                                     |
| `UPV_ACTIVITY_CODE` | Internal activity code (see section below for more info)                                                     |
| `CUSTOM_CRON`       | Allows the user to define a custom [cron trigger](https://github.com/reugn/go-quartz#cron-expression-format) |

</center>

Besides the environment variables, you will have to tell the program when you would like your reservations. The agent will read (and watch for changes of) a file named `config.json`. The file `config.example.json` contains an example of the required content:

```json
{
  "Monday": [
    { "Start": { "Hour": 20 }, "End": { "Hour": 23, "Minute": 59 } },
    { "Start": { "Hour": 18 }, "End": { "Hour": 23, "Minute": 59 } }
  ],
  ...
}
```

Each day of the week (the key) maps to a list of preferred time ranges. Time ranges that come before are assumed to be preferred to time ranges that come later. In this example, the agent would try to make a reservation in a slot in the 20:00-23:59 range. If there are no available slots, it will try in the 18:00-23:59 time range. Valid weekdays are `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Saturday` and `Sunday`

### How to obtain `UPV_ACTIVITY_TYPE` and `UPV_ACTIVITY_CODE`

As the university does not have a public API, these values must be scraped or obtained manually. They are easy enough to obtain manually that there is no need to automate the process.

Go to the [ON-LINE registration of Sportive Activities](https://intranet.upv.es/pls/soalu/sic_depact.HSemActividades) and select your chosen `Propgram` and `Activity` in the website's form. The new URL will contain multiple parameters. Key among them are `p_tipoact` and `p_codacti`, corresponding to `UPV_ACTIVITY_TYPE` and `UPV_ACTIVITY_CODE` respectively. See the image below for a visual explanation.

<center>

<img src="./.img/obtaining-the-codes.png"/>

</center>


## Disclaimer

The UPV **does not authorize** the use of this program. This has been published merely for **educational purposes only**. Also, as there is no official public API, the agent may break at any time and I make no promises on maintaining the project. In case the URL need to be changed, please modify the code in [`lib/requests`](./lib/requests/) and make a pull request.