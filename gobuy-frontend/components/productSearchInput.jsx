"use client"

import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { handleSearch } from "./action";


function SearchInput({currentPage,initialSearch}) {
  const [search, setSearch] = useState(initialSearch);

  return (
    <div className="flex gap-1">
      <Input onChange={(e)=>{setSearch(e.target.value)}} value={search}  placeholder="Search" />
      <Button variant="icon" size='icon' onClick={()=>{handleSearch(currentPage,search)}}>
        <img src="/images/search.png" alt="search"/>
      </Button>
    </div>
  );
}

export default SearchInput;