# TicTacGo

## Крестики-нолики онлайн
-   Backend  in Go (gin-gorm)
-   Frontend in JS (Astro/Svelte)
-   CSS styling using Tailwind
-   Communication using SSE/WS/REST/CRUD
-   Hosting
    -   Frontend -> CloudFlare Pages/Netlify or something similar
	-   Backend  -> Northflank
	-   DataBase -> Planetscale(MySQL)/Upstash(Redis)

## TODO
-	Main page
	-	Rooms list
	-	New room button

-   Room page
    -   Delete room button
	-   Player sign
	-   Current player move
	-   Game field
		-	3x3 table

-	Auth
	-	Cookie based UUID auth
	-	Store cookie using Redis