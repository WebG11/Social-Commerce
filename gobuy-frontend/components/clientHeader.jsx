"use client"
export const Username = ()=>{
    const username = localStorage.getItem("username")
    if (username!==null) {
        username = "";
    }
    return (
        <li>
            {username}
        </li>
    )
}