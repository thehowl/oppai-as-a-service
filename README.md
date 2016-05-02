# oppai-as-a-service

A webapp and API for [oppai](https://github.com/Francesco149/oppai). Because
who the fuck wants to set it up on their own computers when you can just drag
and drop an osr file anyway.

At present, this is somehow close to being released, but there is still need
to do a lot of stuff. Here's what you need to know if you want to try it out:

1. You will need a special version of oppai that I built specifically for OaaS
   and that can be found [here](https://github.com/thehowl/oppai/tree/silent).
   It differs from the original in the fact that the only thing it outputs is
   the amount of total PP. Not even the newline at the end.
2. `mkdir maps`
3. Configuration is done with env variables `DSN` and `OSU_API_KEY`. While
   `OSU_API_KEY` is pretty much self-descriptive, you can find out on how to
   set up the DSN [here](https://github.com/go-sql-driver/mysql#dsn-data-source-name).
   If you're testing locally, `export DSN=root@/oppai` might work.
