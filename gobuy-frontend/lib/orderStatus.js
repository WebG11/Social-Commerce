export const computeStatus = (num)=>{
    if(num === 0){
      return "Pending"
    }else if(num === 1){
      return "Paid"
    }else if(num === 2){
      return "Shipped"
    }else if(num === 3){
      return "Delivered"
    }else if(num ===4){
      return "Cancelled"
    }else if(num ===5){
      return "Refunded"
    }else{
      throw new Error("Invalid status")
    }
  }