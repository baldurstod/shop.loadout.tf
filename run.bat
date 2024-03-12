cls
del .\dist\shop.loadout.tf.exe
go build -ldflags="-X shop.loadout.tf/src/server/server.UseEmbed=false" -o dist/shop.loadout.tf.exe ./src/server/main.go
.\dist\shop.loadout.tf.exe
