package common



// LoadConfig loads the configuration from a YAML or JSON file
func LoadConfig(filePath string) (*Config, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err == nil {
        return &config, nil
    }

    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    return &config, nil
}