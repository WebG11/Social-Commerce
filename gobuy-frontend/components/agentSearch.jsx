"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import Link from "next/link"
import {handleAskAgent} from "./action"

export default ()=> {
  const [isExpanded, setIsExpanded] = useState(false)
  const [query, setQuery] = useState("")
  const [result, setResult] = useState(null)
  const [isLoading, setIsLoading] = useState(false)


  const handleQuery = async () => {
    setResult(null);
    setIsLoading(true)
    try {
      const queryResult = await handleAskAgent(query);
      setResult({ html: queryResult.html })
    } catch (error) {
      console.error("Query failed:", error)
      setResult({ html: "<div>查询失败，请重试。</div>" })
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="fixed bottom-4 right-4">
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="w-16 h-16 rounded-full bg-blue-500 text-white flex items-center justify-center shadow-lg hover:bg-blue-600 transition-colors"
      >
        {isExpanded ? "收起" : "展开"}
      </button>
      {isExpanded && (
        <div className="absolute bottom-20 right-0 w-64 bg-white p-4 rounded-lg shadow-xl">
          <Input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="输入查询内容"
            className="w-full p-2 border border-gray-300 rounded mb-2"
          />
          <Button
            onClick={handleQuery}
            disabled={isLoading}
            className="w-full text-white p-2 rounded transition-colors"
          >
            {isLoading ? "查询中..." : "查询"}
          </Button>
          {result && (
            <div className="mt-4 p-2 bg-gray-100 rounded">
              <div dangerouslySetInnerHTML={{ __html: result.html }} />
              {/* <p>{result.text}</p>
              <div>
                <span>查看商品: </span>
                <ol>
                  <li>
                    <Link href={`/product/${result.productID}`} target="_blank" className="text-blue-500 hover:text-blue-700 underline">{`/product/${result.productID}`}</Link>
                  </li>
                </ol>
              </div> */}
            </div>
          )}
        </div>
      )}
    </div>
  )
}

