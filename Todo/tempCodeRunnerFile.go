 http.HandleFunc("/",func(w http.ResponseWriter, r* http.Request){
        http.ServeFile(w,r,"websocket.html")
    })