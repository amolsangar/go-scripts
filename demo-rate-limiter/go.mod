module example.com/rate-limit

go 1.18

replace example.com/apiConn => ./api

require example.com/apiConn v0.0.0-00010101000000-000000000000

require golang.org/x/time v0.0.0-20220722155302-e5dcc9cfc0b9 // indirect
