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

### Todo

Add progress bar
``` go
count := 10000

tmpl := `{{ bar . "[" "-" ">" " " "]"}} {{percent .}}`
bar := pb.ProgressBarTemplate(tmpl).Start(count)

for i := 0; i < count; i++ {
  bar.Increment()
}
```