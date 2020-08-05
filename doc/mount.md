Separated by whitespace, its 6 columns are:

Mount device if applicable or "none"
Mount point
File system
Mount options
Used by the dump command, 0 to ignore*
Used by the fsck command (which order to check at boot), 0 to ignore*
*Note: mtab places a dummy value into the 5th and 6th columns so that the file retains the same structure as fstab. These columns do not have any meaning in mtab.



https://serverfault.com/questions/267609/how-to-understand-etc-mtab
