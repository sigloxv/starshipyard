# Go Framework 
**Concept** Developing a `cli` framework, a `web` framework, and a `gui`
framework in parallel, we are finding there is quite a bit of underlying code
that makes sense in each codebase. 

Unofficially, this is becomming the application `framework` framework; and so we
have been tracking the differences between these frameworks so we can
consolidate code that makes sense in all three; for example, PID handling,
service management, signal handling, filesystem interaction, and so on. These
are more about interacting with the operating system and so it would make sense
to break off this portion of the code and call it in as a library in all three
frameworks. 

So overtime this needs to be (1) named in such a way that accurately describes
exactly what it does and provides, (2) isolated and made generic enough to work
perfectly in each without any cormprimises, (3) make the code modular as
possible so only the peices needed are called into each of the three frameworks. 


Its either `application framework` or `os library` 
