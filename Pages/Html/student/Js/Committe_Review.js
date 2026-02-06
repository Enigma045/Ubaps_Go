const tabs = document.querySelectorAll(".tabs input");
const tables = document.querySelectorAll(".log-table");

const filter = document.querySelector(".tabs").querySelectorAll("input")
const payload = []




filter.forEach(e =>{
  e.addEventListener("change",()=>{
    payload.length = 0; // Clear the payload array
if (filter[0].checked == true){
   //console.log(filter[0].value)
   payload.push(filter[0].value)
}
if (filter[1].checked == true){
   //console.log(filter[1].value)
   payload.push(filter[1].value)
}
if (filter[2].checked == true){
   //console.log(filter[2].value)
   payload.push(filter[2].value)
}
if (filter[3].checked == true){
   //console.log(filter[3].value)
    payload.push(filter[3].value)
}

console.log(JSON.stringify(payload))

  fetch('/applicants',{
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body:JSON.stringify(payload)
  }).then(response => response.json())
  .then(data => 
    {
      console.log(data)
    }).catch( error => {
    console.error('Error:', error)
  }
    )
  })
})

