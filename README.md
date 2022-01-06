### What is this project for 
    I got kinda sick of bootstraping same configurations for base line services over and over.
    So I created a CLI tool to bootstrap for the necessary configuration files. This project is a *hardcoded* version of that template..
### So what does it contain
    Well, the project contains all the neccessary boiler plate for creating a production ready service. 
    What it currently contains is the following 

    1. We want requests to time out once they reach a configurable threshold.
    2. We want open tracing to some given platform in this instance Jagear and the appropriate context propegation
    3. We want metrics to be collected via prometheus
    4. We want to be able to write routes without getting bogged down with the same boilerplate. 
    5. We want to be able to create workers from some given message bus without the writing the same boilerplate. 


