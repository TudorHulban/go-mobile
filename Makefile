build:
	fyne-cross android --app-id com.example.simpleapp --icon icon.png

emulate:
	@go run -tags mobile .