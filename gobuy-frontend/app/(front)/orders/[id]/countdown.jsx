"use client"
import { useState, useEffect } from 'react';

const CountdownTimer = ({ targetTime }) => {
  const [remainingTime, setRemainingTime] = useState(() => {
    const timestamp = new Date(targetTime).getTime() / 1000 + 10 * 60;
    const now = Date.now() / 1000;
    return timestamp - now;
  });

  useEffect(() => {
    const interval = setInterval(() => {
     setRemainingTime((prevTime => (prevTime > 0 ? prevTime - 1 : 0)));
   }, 1000);

    return () => clearInterval(interval);
  }, [targetTime]);

  const minutes = Math.floor(remainingTime / 60);
  const seconds = parseInt(remainingTime % 60);

  return (
    <>
    {
    remainingTime >0 ?
    <div className='text-lg' suppressHydrationWarning={true}>
        请在: <span className='text-red-500'>{minutes}分{seconds}秒</span> 内付款,否则订单将取消
    </div>:
    <></>
    }
    </>
  );
};

export default CountdownTimer;