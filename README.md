## send commands from the console
![console](https://storage.googleapis.com/rfitzhugh/console.png)

## manage user access
![console](https://storage.googleapis.com/rfitzhugh/user-mgmt.png)

## manage configurations
![console](https://storage.googleapis.com/rfitzhugh/configs.png)

This is mostly a place for me to play around with vue. I repurposed the stuff taught in [this](https://www.udemy.com/course/working-with-vue-3-and-go/) course. That said, I'd still like to post some screenshots and instructions in the event anyone happens upon this repo. Until then, I'm going to assume you know / have a few things.

1. you will need [go](https://go.dev/doc/install) and [node](https://nodejs.org/en/download/)/npm installed.
2. you know how to run code in general (in this case, compiling go and running the vue frontend.)
3. How to install postgres and use the provided brindle.sql (i cant recall why i named it that...) to create the tables you'll need
4. your environment may require additional networking configuration. the backend, frontend, and database need to all run on a separate port and be able to reach one another (in general).
5. you know how to troubleshoot the entire stack (whatever that means)
 
 view run.sh file located in both api/ and frontend/
 
 to create a user with login admin@admin.com and a password of hackerman
 
 ```
 insert into users (email, first_name, last_name, password, user_active, created_at, updated_at)
	values ('admin@admin.com', 'admin', 'admin', '$2a$12$rxOxyleRpZW9y3VPQpXLaeYWgoMy2MU1J7JDWTnN4LpwLsl3tFLhe', 1, '2022-02-01', '2022-02-01')
 ```
 
 do this is after you have provisioned the tables provided by api/data/brindle.sql
 
 [here](https://youtu.be/RML4Hd_oE40) are [some](https://youtu.be/wNGJBhv6F2k) basic videos for getting setup. Please ignore the visual changes in the app you see, I made some significant css/html changes since the video was recorded (that should have no effect on install)
