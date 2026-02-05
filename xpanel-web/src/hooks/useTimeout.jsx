import { useEffect, useRef, useCallback } from 'react';

function useTimeout(callback, delay) {
  const savedCallback = useRef(callback);
  const timeoutRef = useRef();

  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  const clear = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
  }, []);

  useEffect(() => {
    if (typeof delay === 'number') {
      timeoutRef.current = setTimeout(() => savedCallback.current(), delay);
      return clear;
    }
  }, [delay, clear]);

  return { clear }; // 返回一个对象，方便手动控制
}
export default useTimeout