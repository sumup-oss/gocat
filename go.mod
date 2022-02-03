module github.com/sumup-oss/gocat

require (
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/magefile/mage v1.8.0
	github.com/palantir/stacktrace v0.0.0-20161112013806-78658fd2d177
	github.com/spf13/cobra v0.0.3
	github.com/stretchr/testify v1.4.0
	github.com/sumup-oss/go-pkgs v0.0.0-20200306132509-b949afdfe2fe
)

replace (
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20201216223049-8b5274cf687f
	golang.org/x/text => golang.org/x/text v0.3.3
)

go 1.13
