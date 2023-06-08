This is a work in progress.

I don't really want to go down the dot Net or PHP Laravel path and setup database migration scripts just yet

so for now this is just a set of hardcodes SQL schema scripts to:

1) Create the database tables
2) Load some of the data such as the question-form templates and the form questions themselves
3) Create the stored procedures used in this application
   (I know some people hate having any amount of application logic in stored procedures but I find that I am more productive in creating stored procedures and also I think that using stored procedures can help with reducing database server communication round-trips)

