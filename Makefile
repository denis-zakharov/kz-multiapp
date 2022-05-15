generator:
	docker build -t example.com/javagen:v0.1 java-app-generator

transformer:
	docker build -t example.com/java-opts-transformer:v0.1 java-opts-transformer

all: generator transformer