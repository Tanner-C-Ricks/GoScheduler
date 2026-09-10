console.log("Loaded WIndows vars")

window.setUp = function () {

window.DAYS_OF_WEEK = [
                    "Mon",
                    "Tue",
                    "Wed",
                    "Thur",
                    "Fri",
                ];
window.schedule = document.getElementById("scheduleGraph");
window.startHour = 8;
window.endHour = 18;
window.totalMinutesInDay = (window.endHour-window.startHour)*60;
window.increments = 30;
window.incrementsInDay = window.totalMinutesInDay/window.increments;
window.incrementDistance = window.schedule.offsetWidth/window.incrementsInDay; //This would be from 8 to 6 (60 minutes to the hout) and 10 minute increments
window.dayDistance = window.schedule.offsetHeight/(5); //5 days a week
window.blockMin = ((60*3)/window.increments)*window.incrementDistance; //this is a minimum of 30 minutes
window.blockMax = ((60*9)/window.increments)*window.incrementDistance; //this is a maximum of 2hrs

window.maxHours = 20;

console.log("Schedule loaded");
            };