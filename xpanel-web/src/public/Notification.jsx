import useTimeout from '../hooks/useTimeout'

const Notification = ({ props, close }) => {
  // 正确：写在顶层
  // 注意：如果 useTimeout 内部没有处理 props.active，
  // 那么这个 Hook 可能会在组件挂载时立即开始计时
  useTimeout(() => {
    if (props.active) {
      close();
    }
  }, 1500);

  return (
    <>
      {props.active ? (
        <div className={"notification is-light error is-" + props.color} style={{'top': '0px'}}>
          <button className="delete" onClick={close}></button>
          <p>{props.message}</p>
        </div>
      ) : null}
    </>
  )
}

export default Notification