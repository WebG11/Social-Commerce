'use client'
import { Pagination, PaginationContent, PaginationItem, PaginationLink, PaginationNext, PaginationPrevious } from "@/components/ui/pagination"

export default ({totalPage,currentPage,hrefPrefix,hrefParam={}})=>{
  const queryString = '?' + new URLSearchParams(hrefParam).toString();
    return (
        <Pagination className="mt-4">
        <PaginationContent>
          <PaginationItem>
            <PaginationPrevious href={`${hrefPrefix}/${Math.max(currentPage - 1, 1)}${queryString}`} />
          </PaginationItem>
          {Array.from({ length: Math.min(5, Math.ceil(totalPage) - currentPage + 1) }, (_, i) => (
            <PaginationItem key={i}>
              <PaginationLink href={`${hrefPrefix}/${currentPage + i}${queryString}`} isActive={currentPage === currentPage + i}>
                {currentPage + i}
              </PaginationLink>
            </PaginationItem>
          ))}
          <PaginationItem>
            <PaginationNext href={`${hrefPrefix}/${Math.min(currentPage + 1, totalPage)}${queryString}`}/>
          </PaginationItem>
        </PaginationContent>
      </Pagination>
    )
}