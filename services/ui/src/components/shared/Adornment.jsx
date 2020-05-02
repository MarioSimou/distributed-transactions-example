import React from 'react'
import InputAdornment from '@material-ui/core/InputAdornment'

const CustomAdornment = ({position = "start", Icon}) => {
  return (
    <InputAdornment position={position}>
      <Icon/>
    </InputAdornment>
  )
}

export default CustomAdornment
