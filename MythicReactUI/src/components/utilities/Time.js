import {useEffect, useRef } from 'react';

export function toLocalTime(date, view_utc) {
    try {
        if(date === null){
            return "N/A";
        }
        let init_date = new Date(date);
        if (view_utc) {
          return init_date.toDateString() + " " + init_date.toTimeString().substring(0, 8) + " UTC";
        } else {
          let timezoneDate = new Date(date + "Z");
          return timezoneDate.toDateString() + " " + timezoneDate.toLocaleString(['en-us'], {hour12: true, hour: "2-digit", minute: "2-digit"});
        }
        
    } catch (error) {
        console.log("warning", "Failed to get local time converted: " + error.toString());
        return date + " UTC";
    }
}

export function getTimeDifference(checkin_time, current_time) {
    let date = new Date();
    if(current_time !== undefined){
      date = new Date(current_time);
    }
    let now = date.getTime() + date.getTimezoneOffset() * 60000;
    let millisec = Math.abs(now - (new Date(checkin_time).getTime()));
    let seconds = Math.trunc(((millisec / 1000)) % 60);
    let minutes = Math.trunc(((millisec / (1000 * 60))) % 60);
    let hours = Math.trunc(((millisec / (1000 * 60 * 60))) % 24);
    let days = Math.trunc(((millisec / (1000 * 60 * 60 * 24))) % 365);
    let output = "";
    if(days !== 0){ output += days + "d";}
    if(hours !== 0){ output += hours + "h";}
    if(minutes !== 0){ output += minutes + "m";}
    output += seconds + "s";
    return output;
}
export function milisecondsToString(millisec){
    let seconds = Math.trunc(((millisec / 1000)) % 60);
    let minutes = Math.trunc(((millisec / (1000 * 60))) % 60);
    let hours = Math.trunc(((millisec / (1000 * 60 * 60))) % 24);
    let days = Math.trunc(((millisec / (1000 * 60 * 60 * 24))) % 365);
    let output = "";
    if(days !== 0){ output += days + "d";}
    if(hours !== 0){ output += hours + "h";}
    if(minutes !== 0){ output += minutes + "m";}
    output += seconds + "s";
    return output;
}
//https://overreacted.io/making-setinterval-declarative-with-react-hooks/
export function useInterval(callback, delay, mountedRef, parentMountedRef) {
  const savedCallback = useRef();

  useEffect(() => {
    savedCallback.current = callback;
  });
  useEffect(() => {
    function tick() {
      if((mountedRef && !mountedRef.current) || (parentMountedRef && !parentMountedRef.current)){
        //console.log("returning")
        return;
      }
      savedCallback.current();
    }
    if( (mountedRef && !mountedRef.current) || (parentMountedRef && !parentMountedRef.current)){
      //console.log("returning2")
      return;
    }
    let id = setInterval(tick, delay);
    return () => clearInterval(id);
  }, [delay, mountedRef, parentMountedRef]);
}
