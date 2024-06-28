FROM golang:1.21

WORKDIR /tnq/apps/xml-central/xmlc-wrapper-go

COPY . /tnq/apps/xml-central/xmlc-wrapper-go

EXPOSE 8081

CMD ["./main"]