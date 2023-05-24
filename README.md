# dg
A fast data generator that produces CSV files from generated relational data

### Concepts

dg takes its configuration from a config file that is parsed in the form of an array of objects. Each object represents a CSV file to be generated for a named table and contains a collection of columns to generate data for.

There are three ways to generate data for columns:

`gen`  - Generate a random value for the column. Here's an example:

``` yaml
- name: id
  type: gen
  processor:
    value: uuid
```

This configuration will generate a random UUID for the id column. `value` points to a function in the code that generates the data. A list of available functions can be found [here](https://github.com/codingconcepts/dg#functions).

**`ref`**  - References a value from a previously generated table. Here's an example:

``` yaml
- name: ptype
  type: ref
  processor:
    table: person_type
    column: id
```

This configuration will choose a random id from the person_type table and create a **`ptype`** column to store the values.

Use the `ref` type if you need to reference another table but don't need to generate a new row for *every* instance of the referenced column.

**`each`** - Creates a row for each value in another table. If multiple `each` columns are provided, a Cartesian product of both columns will be generated.

Here's an example of one `each` column:

``` yaml
- table: person
  count: 3
  columns:
    - name: id
      type: gen
      processor:
        value: uuid

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
        value: uuid

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
        value: uuid

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

### Usage

```
$ dg
Usage dg:
  -c string
        the absolute or relative path to the config file
  -o string
        the absolute or relative path to the output dir (default ".")
```

Create a config file. In the following example, we're creating 10,000 people, 50 events, 5 person types, and then populating the many-to-many `person_event` resolver table:

``` yaml
- table: person
  count: 10000
  columns:
    - name: id
      type: gen
      processor:
        value: uuid

- table: event
  count: 50
  columns:
    - name: id
      type: gen
      processor:
        value: uuid

- table: person_type
  count: 5
  columns:
    - name: id
      type: gen
      processor:
        value: uuid
    - name: name
      type: gen
      processor:
        value: int8

- table: person_event
  columns:
    - name: id
      type: gen
      processor:
        value: uuid
    - name: person_type
      type: ref
      processor:
        table: person_type
        column: id
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
    
```

Run the application:
```
$ dg -c your_config_file.yaml -o your_output_dir 
```

### Functions

| Name | Example |
| ---- | ------- |
| address | 1901 North Midwest Boulevard |
| amount_with_currency | XOF 5151.600000 |
| bool | true |
| cc_number | 6011899034217690 |
| cc_type | Discover |
| century | XVIII |
| city | Arvada |
| currency | UZS |
| date | 2014-02-07 |
| day_of_month |  2 |
| day_of_week | Friday |
| domain_name | lKZVgGU.info |
| e164_phone_number | +719653721048 |
| email | kOgopwo@TGxcrij.net |
| first_name | Adela |
| first_name_female | Ruth |
| first_name_male | Nash |
| int16 | 4073 |
| int32 | 1567839306 |
| int64 | 3223704240700754243 |
| int8 | 114 |
| ipv4 | 9.31.166.108 |
| ipv6 | 3b10:66f8:fa46:ee01:c520:2a63:d910:7ff9 |
| last_name | Kessler |
| latitude | -39.808380126953125 |
| longitude | -81.36772155761719 |
| mac_address | d3:4a:49:fe:51:4a |
| month_name | February |
| name | Prof. Hardy Howell |
| paragraph | Ad sapiente pariatur fugit soluta omnis. |
| password | eBPyDCVsiPrXadlSCIbmcgMWesKcgooyvocwKgrIuWStRfdKPM |
| phone_number | 623-107-1458 |
| postcode | 72701 |
| sentence | Qui mollitia qui totam fuga qui. |
| state | TN |
| time | 12:59:08 |
| timeperiod | AM |
| timestamp | 1980-08-21 13:30:30 |
| timezone | Europe/Podgorica |
| title_female | Dr. |
| title_male | Mr. |
| toll_free_phone_number | (888) 156-793824 |
| unix_time | 169644490 |
| url | http://erMKGUs.net/grVsTVn.html |
| user_name | qAgFMOv |
| uuid | 6fe8ac8ea3a246b5bfafe1da24bf084f |
| uuid_hyphen | 05acae1b-3db8-40b4-bd70-287b55d5e026 |
| word | quo |
| year | 2018 |

### Todos

Add progress bar
``` go
count := 10000

tmpl := `{{ bar . "[" "-" ">" " " "]"}} {{percent .}}`
bar := pb.ProgressBarTemplate(tmpl).Start(count)

for i := 0; i < count; i++ {
  bar.Increment()
}
```