import { useEffect } from 'react'
const Notification = ({ props, close }) => {
  useEffect(() => {
    setTimeout(function(){
      close()
    },1500)
  })
  return (
    <>
      {props.active ?
        <div className={"notification is-light error is-" + props.color} style={{'top': '0px'}}>
          <button className="delete" onClick={close}></button>
          <p>{props.message}</p>
        </div>
        :
        <></>
      }
    </>
  )
}
export default Notification