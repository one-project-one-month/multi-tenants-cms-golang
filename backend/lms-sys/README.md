

####

```go

    databaseCtx, cancel := cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel
    
    
```