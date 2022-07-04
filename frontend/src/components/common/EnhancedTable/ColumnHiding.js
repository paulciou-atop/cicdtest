import React from 'react'
import FormGroup from "@mui/material/FormGroup";
import FormControlLabel from "@mui/material/FormControlLabel";
import Checkbox from "@mui/material/Checkbox";

export const ColumnHiding = ({column,isSelected,handleClick,}) => {

  return (
    <>
    {
        column.map((col)=>{
            const isItemSelected = isSelected(col["key"])
            return  (  
            <FormGroup title='col hide' key={col.key}>
                <FormControlLabel
                key={col.key}
                label={col.label}
                control={<Checkbox 
                 checked={isItemSelected}
                // value={col.id}
                   value={col.key}
                 onChange={(event) => handleClick(event, col["key"])}
                />}
                 /> 
            </FormGroup>
            )
        })
    }    
    </>
  )
}
