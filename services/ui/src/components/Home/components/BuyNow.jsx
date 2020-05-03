import React from 'react'
import {
  Fab,
  Typography,
} from '@material-ui/core'
import PaymentIcon from '@material-ui/icons/Payment';

const BuyNow = ({ onClickBuyNow, ...other}) => {
  return (
    <Typography component="div" {...other}>
      <Fab color="primary" size="medium" variant="extended" onClick={onClickBuyNow}>
        <PaymentIcon/>&nbsp;Buy
      </Fab>
    </Typography>
  )
}

export default BuyNow