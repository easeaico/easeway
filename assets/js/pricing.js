"use strict";

/* ======= Popover ======= */
var popoverTriggerList = [].slice.call(document.querySelectorAll('.info-popover-trigger'))
var popoverList = popoverTriggerList.map(function (popoverTriggerEl) {
  return new bootstrap.Popover(popoverTriggerEl)
})



/* ======= responsive pricing table ======== */
const pricingTabs = document.querySelectorAll('#pricing-tabs .pricing-tab');
const pricingTable = document.querySelector('#pricing-table');
const pricingDataCells = pricingTable.querySelectorAll('.pricing-data');


pricingTabs.forEach((pricingTab) => {
	
  pricingTab.addEventListener('click', () => {
	  
	    // Highlight selected pricing plan tab 
	    
	    /* ref:  for...of loop - https://stackoverflow.com/questions/48962644/implement-siblings-in-vanilla-javascript */
	    
	    for (let siblingTab of pricingTab.parentNode.children) {
	        siblingTab.classList.remove('active');
	    }
	    
	    pricingTab.classList.add('active');
	      
	    
	    //Show selected pricing plan table content
	    
	    let dataTarget = pricingTab.getAttribute('data-target');
	    
	    //console.log(dataTarget);
	    
	    for (let pricingDataCell of pricingDataCells) {
	        pricingDataCell.style.display = "none";
	    }
	    
	    for (let dataTargetCell of pricingTable.querySelectorAll('.' + dataTarget) ) {
		    dataTargetCell.style.display = "table-cell";
	    }

    
    });
});