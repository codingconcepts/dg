<p align="center">
  <img src="assets/cover.png" alt="drawing" width="150"/>
</p>

A fast data generator that produces CSV files from generated relational data

## Table of Contents
1. [Installation](#installation)
1. [Usage](#usage)
1. [Processors](#processors)
    * [gen](#gen)
    * [set](#set)
    * [inc](#inc)
    * [ref](#ref)
    * [each](#each)
    * [range](#range)
    * [match](#match)
1. [Inputs](#inputs)
    * [csv](#csv)
1. [Functions](#functions)
1. [Thanks](#thanks)
1. [Todos](#todos)

### Installation

Find the release that matches your architecture on the [releases](https://github.com/codingconcepts/dg/releases) page.

Download the tar, extract the executable, and move it into your PATH:

```
$ tar -xvf dg_[VERSION]-rc1_macOS.tar.gz
```

### Usage

```
$ dg
Usage dg:
  -c string
        the absolute or relative path to the config file
  -o string
        the absolute or relative path to the output dir (default ".")
```

Create a config file. In the following example, we create 10,000 people, 50 events, 5 person types, and then populate the many-to-many `person_event` resolver table with 500,000 rows that represent the Cartesian product between the person and event tables:

``` yaml
tables:
  - table: person
    count: 10000
    columns:
      # Generate a random UUID for each person
      - name: id
        type: gen
        processor:
          value: ${uuid}

  - table: event
    count: 50
    columns:
      # Generate a random UUID for each event
      - name: id
        type: gen
        processor:
          value: ${uuid}

  - table: person_type
    count: 5
    columns:
      # Generate a random UUID for each person_type
      - name: id
        type: gen
        processor:
          value: ${uuid}
      
      # Generate a random 16 bit number and left-pad it to 5 digits
      - name: name
        type: gen
        processor:
          value: ${uint16}
          format: "%05d"

  - table: person_event
    columns:
      # Generate a random UUID for each person_event
      - name: id
        type: gen
        processor:
          value: ${uuid}
      
      # Select a random id from the person_type table
      - name: person_type
        type: ref
        processor:
          table: person_type
          column: id
      
      # Generate a person_id column for each id in the person table
      - name: person_id
        type: each
        processor:
          table: person
          column: id
      
      # Generate an event_id column for each id in the event table
      - name: event_id
        type: each
        processor:
          table: event
          column: id
```

Run the application:
```
$ dg -c your_config_file.yaml -o your_output_dir
loaded config file                       took: 428µs
generated table: person                  took: 41ms
generated table: event                   took: 159µs
generated table: person_type             took: 42µs
generated table: person_event            took: 1s
generated all tables                     took: 1s
wrote csv: person                        took: 1ms
wrote csv: event                         took: 139µs
wrote csv: person_type                   took: 110µs
wrote csv: person_event                  took: 144ms
wrote all csvs                           took: 145ms
```

This will output:

```
your_output_dir
├── event.csv
├── person.csv
├── person_event.csv
└── person_type.csv
```

If you're following along locally, spin up a local web server using something like python's `http.server`:

```
$ python3 -m http.server 3000 -d your_output_dir
```

Then import the files as you would any other; here's an example insert into CockroachDB:

``` sql
IMPORT INTO "person" ("id")
CSV DATA (
    'http://localhost:3000/person.csv'
)
WITH skip='1', nullif = '', allow_quoted_null;

IMPORT INTO "event" ("id")
CSV DATA (
    'http://localhost:3000/event.csv'
)
WITH skip='1', nullif = '', allow_quoted_null;

IMPORT INTO "person_type" ("id", "name")
CSV DATA (
    'http://localhost:3000/person_type.csv'
)
WITH skip='1', nullif = '', allow_quoted_null;

IMPORT INTO "person_event" ("person_id", "event_id", "id", "person_type")
CSV DATA (
    'http://localhost:3000/person_event.csv'
)
WITH skip='1', nullif = '', allow_quoted_null;
```

### Processors

dg takes its configuration from a config file that is parsed in the form of an object containing arrays of objects; `tables` and `inputs`. Each object in the `tables` array represents a CSV file to be generated for a named table and contains a collection of columns to generate data for.

##### gen

Generate a random value for the column. Here's an example:

``` yaml
- name: sku
  type: gen
  processor:
    value: SKU${uint16}
    format: "%05d"
```

This configuration will generate a random left-padded `uint16` with a prefix of "SKU" for a column called "sku". `value` contains zero or more function placeholders that can be used to generate data. A list of available functions can be found [here](https://github.com/codingconcepts/dg#functions).

##### set

Select a value from a given set. Here's an example:

``` yaml
- name: user_type
  type: set
  processor:
    values: [admin, regular, read-only]
```

This configuration will select between the values "admin", "regular", and "read-only"; each with an equal probability of being selected.

Items in a set can also be given a weight, which will affect their likelihood of being selected. Here's an example:

``` yaml
- name: favourite_animal
  type: set
  processor:
    values: [rabbit, dog, cat]
    weights: [10, 60, 30]
```

This configuration will select between the values "rabbit", "dog", and "cat"; each with different probabilities of being selected. Rabbits will be selected approximately 10% of the time, dogs 60%, and cats 30%. The total value doesn't have to be 100, however, you can use whichever numbers make most sense to you.

##### inc

Generates an incrementing number. Here's an example:

``` yaml
- name: id
  type: inc
  processor:
    start: 1
    format: "P%03d"
```

This configuration will generate left-padded ids starting from 1, and format them with a prefix of "P".

##### ref

References a value from a previously generated table. Here's an example:

``` yaml
- name: ptype
  type: ref
  processor:
    table: person_type
    column: id
```

This configuration will choose a random id from the person_type table and create a `ptype` column to store the values.

Use the `ref` type if you need to reference another table but don't need to generate a new row for *every* instance of the referenced column.

##### each

Creates a row for each value in another table. If multiple `each` columns are provided, a Cartesian product of both columns will be generated.

Here's an example of one `each` column:

``` yaml
- table: person
  count: 3
  columns:
    - name: id
      type: gen
      processor:
        value: ${uuid}

# person
#
# id
# c40819f8-2c76-44dd-8c44-5eef6a0f2695
# 58f42be2-6cc9-4a8c-b702-c72ab1decfea
# ccbc2244-667b-4bb5-a5cd-a1e9626a90f9

- table: pet
  columns:
    - name: person_id
      type: each
      processor:
        table: person
        column: id
    - name: name
      type: gen
      processor:
        value: first_name

# pet
#
# person_id                            name
# c40819f8-2c76-44dd-8c44-5eef6a0f2695 Carlo
# 58f42be2-6cc9-4a8c-b702-c72ab1decfea Armando
# ccbc2244-667b-4bb5-a5cd-a1e9626a90f9 Kailey
```

Here's an example of two `each` columns:

``` yaml
- table: person
  count: 3
  columns:
    - name: id
      type: gen
      processor:
        value: ${uuid}

# person
#
# id
# c40819f8-2c76-44dd-8c44-5eef6a0f2695
# 58f42be2-6cc9-4a8c-b702-c72ab1decfea
# ccbc2244-667b-4bb5-a5cd-a1e9626a90f9

- table: event
  count: 3
  columns:
    - name: id
      type: gen
      processor:
        value: ${uuid}

# event
#
# id
# 39faeb54-67d1-46db-a38b-825b41bfe919
# 7be981a9-679b-432a-8a0f-4a0267170c68
# 9954f321-8040-4cd7-96e6-248d03ee9266

- table: person_event
  columns:
    - name: person_id
      type: each
      processor:
        table: person
        column: id
    - name: event_id
      type: each
      processor:
        table: event
        column: id

# person_event
#
# person_id                            
# c40819f8-2c76-44dd-8c44-5eef6a0f2695 39faeb54-67d1-46db-a38b-825b41bfe919
# c40819f8-2c76-44dd-8c44-5eef6a0f2695 7be981a9-679b-432a-8a0f-4a0267170c68
# c40819f8-2c76-44dd-8c44-5eef6a0f2695 9954f321-8040-4cd7-96e6-248d03ee9266
# 58f42be2-6cc9-4a8c-b702-c72ab1decfea 39faeb54-67d1-46db-a38b-825b41bfe919
# 58f42be2-6cc9-4a8c-b702-c72ab1decfea 7be981a9-679b-432a-8a0f-4a0267170c68
# 58f42be2-6cc9-4a8c-b702-c72ab1decfea 9954f321-8040-4cd7-96e6-248d03ee9266
# ccbc2244-667b-4bb5-a5cd-a1e9626a90f9 39faeb54-67d1-46db-a38b-825b41bfe919
# ccbc2244-667b-4bb5-a5cd-a1e9626a90f9 7be981a9-679b-432a-8a0f-4a0267170c68
# ccbc2244-667b-4bb5-a5cd-a1e9626a90f9 9954f321-8040-4cd7-96e6-248d03ee9266
```

Use the `each` type if you need to reference another table and need to generate a new row for *every* instance of the referenced column.

##### range

Generates data within a given range. Note that a number of factors determine how this generator will behave. The step (and hence, number of rows) will be generated in the following priority order:

1. If an `each` generator is being used, step will be derived from that
1. If a `count` is provided, step will be derived from that
1. Otherwise, `step` will be used

Here's an example that generates all dates between `2020-01-01` and `2023-01-01` at daily intervals:

``` yaml
- table: event
  columns:
    - name: date
      type: range
      processor:
        type: date
        from: 2020-01-01
        to: 2023-01-01
        step: 24h
        format: 2006-01-02
```

Here's an example that generates 10 dates between `2020-01-01` and `2023-01-02`:

``` yaml
- table: event
  count: 10
  columns:
    - name: date
      type: range
      processor:
        type: date
        from: 2020-01-01
        to: 2023-01-01
        format: 2006-01-02
        step: 24h   # Ignored due to table count.
```

Here's an example that generates 20 dates (one for every row found from an `each` generator) between `2020-01-01` and `2023-01-02`:

``` yaml
- table: person
  count: 20
  columns:
    - name: id
      type: gen
      processor:
        value: ${uuid}

- table: event
  count: 10         # Ignored due to resulting count from "each" generator.
  columns:
    - name: person_id
      type: each
      processor:
        table: person
        column: id

    - name: date
      type: range
      processor:
        type: date
        from: 2020-01-01
        to: 2023-01-01
        format: 2006-01-02
```

##### match

Generates data by matching data in another table. In this example, we'll assume there's a CSV file for the `significant_event` input that generates the following table:

|date|event|
|----|-----|
|2023-01-10|abc|
|2023-01-11||
|2023-01-12|def|

``` yaml
inputs:
  - name: significant_event
    type: csv
    source:
      file_name: significant_dates.csv

tables:
  - table: events
    columns:
      - name: timeline_date
        type: range
        processor:
          type: date
          from: 2023-01-09
          to: 2023-01-13
          format: 2006-01-02
          step: 24h
      - name: timeline_event
        type: match
        processor:
          source_table: significant_event
          source_column: date
          source_value: events
          match_column: timeline_date
```

dg will match rows in the significant_event table with rows in the event column based on the match between `significant_event.date` and `events.date`, and take the value from the `significant_events.events` column where there's a match (otherwise leaving `NULL`). This will result in the following `events` table being generated:

|timeline_date|timeline_event|
|----|-----|
|2023-01-09||
|2023-01-10|abc|
|2023-01-11||
|2023-01-12|def|
|2023-01-13||

### Inputs

dg takes its configuration from a config file that is parsed in the form of an object containing arrays of objects; `tables` and `inputs`. Each object in the `inputs` array represents a data source from which a table can be created. Tables created via inputs will not result in output CSVs.

##### csv

Reads in a CSV file as a table that can be referenced from other tables. Here's an example:

``` yaml
- name: significant_event
  type: csv
  source:
    file_name: significant_dates.csv
```

This configuration will read from a file called significant_dates.csv and create a table from its contents. Note that the `file_name` should be relative to the config directory, so if your CSV file is in the same directory as your config file, just include the file name.


### Functions

| Name | Type | Example |
| ---- | ---- | ------- |
| ${ach_account} | string | 586981797546 |
| ${ach_routing} | string | 441478502 |
| ${adjective_demonstrative} | string | there |
| ${adjective_descriptive} | string | eager |
| ${adjective_indefinite} | string | several |
| ${adjective_interrogative} | string | whose |
| ${adjective_possessive} | string | her |
| ${adjective_proper} | string | Iraqi |
| ${adjective_quantitative} | string | sufficient |
| ${adjective} | string | double |
| ${adverb_degree} | string | far |
| ${adverb_frequency_definite} | string | daily |
| ${adverb_frequency_indefinite} | string | always |
| ${adverb_manner} | string | unexpectedly |
| ${adverb_place} | string | here |
| ${adverb_time_definite} | string | yesterday |
| ${adverb_time_indefinite} | string | just |
| ${adverb} | string | far |
| ${animal_type} | string | mammals |
| ${animal} | string | ape |
| ${app_author} | string | RedLaser |
| ${app_name} | string | SlateBlueweek |
| ${app_version} | string | 3.2.10 |
| ${bitcoin_address} | string | 16YmZ5ol5aXKjilZT2c2nIeHpbq |
| ${bitcoin_private_key} | string | 5JzwyfrpHRoiA59Y1Pd9yLq52cQrAXxSNK4QrGrRUxkak5Howhe |
| ${bool} | bool | true |
| ${breakfast} | string | Awesome orange chocolate muffins |
| ${bs} | string | leading-edge |
| ${car_fuel_type} | string | LPG |
| ${car_maker} | string | Seat |
| ${car_model} | string | Camry Solara Convertible |
| ${car_transmission_type} | string | Manual |
| ${car_type} | string | Passenger car mini |
| ${chrome_user_agent} | string | Mozilla/5.0 (X11; Linux i686) AppleWebKit/5310 (KHTML, like Gecko) Chrome/37.0.882.0 Mobile Safari/5310 |
| ${city} | string | Memphis |
| ${color} | string | DarkBlue |
| ${company_suffix} | string | LLC |
| ${company} | string | PlanetEcosystems |
| ${connective_casual} | string | an effect of |
| ${connective_complaint} | string | i.e. |
| ${connective_examplify} | string | for example |
| ${connective_listing} | string | next |
| ${connective_time} | string | soon |
| ${connective} | string | for instance |
| ${country_abr} | string | VU |
| ${country} | string | Eswatini |
| ${credit_card_cvv} | string | 315 |
| ${credit_card_exp} | string | 06/28 |
| ${credit_card_type} | string | Mastercard |
| ${currency_long} | string | Mozambique Metical |
| ${currency_short} | string | SCR |
| ${date} | time.Time | 2005-01-25 22:17:55.371781952 +0000 UTC |
| ${day} | int | 27 |
| ${dessert} | string | Chocolate coconut dream bars |
| ${dinner} | string | Creole potato salad |
| ${domain_name} | string | centralb2c.net |
| ${domain_suffix} | string | com |
| ${email} | string | ethanlebsack@lynch.name |
| ${emoji} | string | ♻️ |
| ${file_extension} | string | csv |
| ${file_mime_type} | string | image/vasa |
| ${firefox_user_agent} | string | Mozilla/5.0 (X11; Linux x86_64; rv:6.0) Gecko/1951-07-21 Firefox/37.0 |
| ${first_name} | string | Kailee |
| ${flipacoin} | string | Tails |
| ${float32} | float32 | 2.7906555e+38 |
| ${float64} | float64 | 4.314310154193861e+307 |
| ${fruit} | string | Eggplant |
| ${gender} | string | female |
| ${hexcolor} | string | #6daf06 |
| ${hobby} | string | Bowling |
| ${hour} | int | 18 |
| ${http_method} | string | DELETE |
| ${http_status_code_simple} | int | 404 |
| ${http_status_code} | int | 503 |
| ${http_version} | string | HTTP/1.1 |
| ${int16} | int16 | 18940 |
| ${int32} | int32 | 2129368442 |
| ${int64} | int64 | 5051946056392951363 |
| ${int8} | int8 | 110 |
| ${ipv4_address} | string | 191.131.155.85 |
| ${ipv6_address} | string | 1642:94b:52d8:3a4e:38bc:4d87:846e:9c83 |
| ${job_descriptor} | string | Senior |
| ${job_level} | string | Identity |
| ${job_title} | string | Executive |
| ${language_abbreviation} | string | kn |
| ${language} | string | Bengali |
| ${last_name} | string | Friesen |
| ${latitude} | float64 | 45.919913 |
| ${longitude} | float64 | -110.313125 |
| ${lunch} | string | Sweet and sour pork balls |
| ${mac_address} | string | bd:e8:ce:66:da:5b |
| ${minute} | int | 23 |
| ${month_string} | string | April |
| ${month} | int | 10 |
| ${name_prefix} | string | Ms. |
| ${name_suffix} | string | I |
| ${name} | string | Paxton Schumm |
| ${nanosecond} | int | 349669923 |
| ${nicecolors} | []string | [#490a3d #bd1550 #e97f02 #f8ca00 #8a9b0f] |
| ${noun_abstract} | string | timing |
| ${noun_collective_animal} | string | brace |
| ${noun_collective_people} | string | mob |
| ${noun_collective_thing} | string | orchard |
| ${noun_common} | string | problem |
| ${noun_concrete} | string | town |
| ${noun_countable} | string | cat |
| ${noun_uncountable} | string | wisdom |
| ${noun} | string | case |
| ${opera_user_agent} | string | Opera/10.10 (Windows NT 5.01; en-US) Presto/2.11.165 Version/13.00 |
| ${password} | string | 1k0vWN 9Z|4f={B YPRda4ys. |
| ${pet_name} | string | Bernadette |
| ${phone_formatted} | string | (476)455-2253 |
| ${phone} | string | 2692528685 |
| ${phrase} | string | I'm straight |
| ${preposition_compound} | string | ahead of |
| ${preposition_double} | string | next to |
| ${preposition_simple} | string | at |
| ${preposition} | string | outside of |
| ${programming_language} | string | PL/SQL |
| ${pronoun_demonstrative} | string | those |
| ${pronoun_interrogative} | string | whom |
| ${pronoun_object} | string | us |
| ${pronoun_personal} | string | I |
| ${pronoun_possessive} | string | mine |
| ${pronoun_reflective} | string | yourself |
| ${pronoun_relative} | string | whom |
| ${pronoun} | string | those |
| ${quote} | string | "Raw denim tilde cronut mlkshk photo booth kickstarter." - Gunnar Rice |
| ${rgbcolor} | []int | [152 74 172] |
| ${safari_user_agent} | string | Mozilla/5.0 (Windows; U; Windows 95) AppleWebKit/536.41.5 (KHTML, like Gecko) Version/5.2 Safari/536.41.5 |
| ${safecolor} | string | gray |
| ${second} | int | 58 |
| ${snack} | string | Crispy fried chicken spring rolls |
| ${ssn} | string | 783135577 |
| ${state_abr} | string | AL |
| ${state} | string | Kentucky |
| ${street_name} | string | Way |
| ${street_number} | string | 6234 |
| ${street_prefix} | string | Port |
| ${street_suffix} | string | stad |
| ${street} | string | 11083 Lake Fall mouth |
| ${time_zone_abv} | string | ADT |
| ${time_zone_full} | string | (UTC-02:00) Coordinated Universal Time-02 |
| ${time_zone_offset} | float32 | 3 |
| ${time_zone_region} | string | Asia/Aqtau |
| ${time_zone} | string | Mountain Standard Time (Mexico) |
| ${uint128_hex} | string | 0xcd50930d5bc0f2e8fa36205e3d7bd7b2 |
| ${uint16_hex} | string | 0x7c80 |
| ${uint16} | uint16 | 25076 |
| ${uint256_hex} | string | 0x61334b8c51fa841bf9a3f1f0ac3750cd1b51ca2046b0fb75627ac73001f0c5aa |
| ${uint32_hex} | string | 0xfe208664 |
| ${uint32} | uint32 | 783098878 |
| ${uint64_hex} | string | 0xc8b91dc44e631956 |
| ${uint64} | uint64 | 5722659847801560283 |
| ${uint8_hex} | string | 0x65 |
| ${uint8} | uint8 | 192 |
| ${url} | string | https://www.leadcutting-edge.net/productize |
| ${user_agent} | string | Opera/10.64 (Windows NT 5.2; en-US) Presto/2.13.295 Version/10.00 |
| ${username} | string | Gutmann2845 |
| ${uuid} | string | e6e34ff4-1def-41e5-9afb-f697a51c0359 |
| ${vegetable} | string | Tomato |
| ${verb_action} | string | knit |
| ${verb_helping} | string | did |
| ${verb_linking} | string | has |
| ${verb} | string | be |
| ${weekday} | string | Tuesday |
| ${word} | string | month |
| ${year} | int | 1962 |
| ${zip} | string | 45618 |

### Building releases locally

```
$ VERSION=0.1.0 make release
```

### Thanks

Thanks to the maintainers of the following fantastic packages, whose code this tools makes use of:

* [samber/lo](https://github.com/samber/lo)
* [brianvoe/gofakeit](https://github.com/brianvoe/gofakeit)
* [go-yaml/yaml](https://github.com/go-yaml/yaml)
* [stretchr/testify](github.com/stretchr/testify/assert)

### Todos

* Write file after generating, then only keep columns that other tables need

* Add progress bar
``` go
count := 10000

tmpl := `{{ bar . "[" "-" ">" " " "]"}} {{percent .}}`
bar := pb.ProgressBarTemplate(tmpl).Start(count)

for i := 0; i < count; i++ {
  bar.Increment()
}
```
