# dg
A fast data generator that produces CSV files from generated relational data

### Usage

Create a config file:
``` yaml
- table: person
  count: 100000
  columns:
    - name: id
      value: ${gen uuid}
    - name: full_name
      value: ${gen first_name} ${gen last_name}
    - name: dob
      value: ${gen date}
```

Run the application:
```
$ dg -c your_config_file.yaml -o your_output_dir 
```

### Todos

Support for multiple foreach tables, resulting in a Cartesian product of both:
``` yaml
- table: person
  ...

- table: event
  ...

- table: person_event
  foreach: [person, event]
```

For example:
```
# person
id,name
1,a
2,b

# event
id,name
3,wacken
4,download

# person_event
person_id,event_id
1,3
1,4
2,3
2,4
```

Add progress bar
``` go
count := 10000

tmpl := `{{ bar . "[" "-" ">" " " "]"}} {{percent .}}`
bar := pb.ProgressBarTemplate(tmpl).Start(count)

for i := 0; i < count; i++ {
  bar.Increment()
}
```