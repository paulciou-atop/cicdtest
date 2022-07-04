import React,{useState} from 'react'
import MUITable from "../../components/common/EnhancedTable/MUITable"
import Popup from "../../components/common/Popup";
import ViewColumnIcon from '@mui/icons-material/ViewColumn';
import { ColumnHiding } from '../../components/common/EnhancedTable/ColumnHiding';
import inventoryData from '../../utils/data/inventoryData.json'

const column = [
  {  key: "deviceType", label: "Device Type" },
  {  key: "model", label: "Model" },
  {  key: "ipAddress", label: "IP Address" },
  {  key: "macAddress", label: "Mac Address" },
  {  key: "hostName", label: "Host Name"},
  {  key: "kernel",label:"Kernel"},
  {  key:"ap",label:"Ap"},
  {  key:"createdAt",label:"Created At"},
  {  key:"lastSeen",label:"Last Seen"},
  {  key:"lastMissing",label:"Last Missinng"},
  {  key:"status",label:"Status",disableSort: true}
];

export const InventoryDevice = () => {
 const [openPopup, setopenPopup] = useState(false);
 const [selectedCol,setSelectedCol]=useState(column.map(col=>col.key));
 const isSelected=(id)=>selectedCol.indexOf(id) !== -1;
 const handleClick = (event, id) => {
        const selectedIndex = selectedCol.indexOf(id);    
        let newSelected = [];
        if (selectedIndex === -1) {
          newSelected = newSelected.concat(selectedCol, id);
        } else if (selectedIndex === 0) {
          newSelected = newSelected.concat(selectedCol.slice(1));
        } else if (selectedIndex === selectedCol.length - 1) {
          newSelected = newSelected.concat(selectedCol.slice(0, -1));
        } else if (selectedIndex > 0) {
          newSelected = newSelected.concat(
            selectedCol.slice(0, selectedIndex),
            selectedCol.slice(selectedIndex + 1)
          );
        }
        setSelectedCol(newSelected)
 };

  return (
    <>
    <MUITable 
    DataSource={inventoryData} 
    Columns={column}
    selectedCol={selectedCol}
    rowKey="id"
    title={`Inventory Management`}
    TotalRowCount={inventoryData.length}
    currentPage={1}
    options={{ sortable: true, selectable: false, filtrable: true }}
    toolbarOptions={{
        exportData: true,
        globalFilter: true,
        columnSelected:{
        enable:true,
        label:"columns",
        icon: <ViewColumnIcon/>,
        color:"primary",
        onClick: (event) => {
          setopenPopup(true);
        },
      }
    }}
    />
     <Popup
        title="Columns Hiding"
        openPopup={openPopup}
        setOpenPopup={setopenPopup}
      >
      <ColumnHiding 
      column={column} 
      handleClick={handleClick} 
      isSelected={isSelected} 
      />
      </Popup>
   
    </>
  )
}
