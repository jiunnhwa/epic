package worker

// func GetNews(urls []string, interval time.Duration) {
// 	NextURL := func(list []string) func() string {
// 		len := len(list)
// 		currNum := -1
// 		return func() string {
// 			currNum += 1
// 			if currNum >= len {
// 				currNum = 0
// 			}
// 			return list[currNum]
// 		}
// 	}(urls)

// 	//url := (GetNextURL(URLs))
// 	ticker := time.NewTicker(interval)
// 	for ; true; <-ticker.C {
// 		u := NextURL() //url()
// 		fname := path.Join("inbox", GetHostName(u)+"-"+"data.txt")
// 		fmt.Println(fname, u)
// 		//WriteData(fname, (Fetch("GET", u, "")))
// 	}
// }

// func NewsAPI() {
// 	//job := NewJob(r.Method + "/" + request.GetAction(r) + "/" + courseid)
// 	//mylogger.Trace.Println("START:", job.ID)
// 	defer func() {
// 		//	mylogger.Trace.Println("END:", job.ID)
// 		//	job.End()
// 		//	job = nil
// 	}()
// }
