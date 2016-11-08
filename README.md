# cron
Based on github.com/robfig/cron

### Example
```Go
package main
                               
import (
        "fmt"
        "github.com/mmadfox/cron"       
) 

func main() {                  
        c, err := cron.New("@every 1m") 
        if err != nil {        
                panic(err)     
        }
        c.Handle(func() error {
                fmt.Println("Run handle ok")    
                return nil     
        })
        c.CloseHandle(func() error {    
                fmt.Println("Close handle ok") 
                return nil     
        })
        c.Run()                
}   
```
