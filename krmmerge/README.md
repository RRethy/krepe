The 3-way merge algorithm operates both on the level of each resource and on each individual field with a resource.

On the resource level, the rules are:

A resource present in none will not be present in local.
A resource missing from origin and local but present in upstream will be added to local.
A resource present in local but missing from origin and upstream will be kept without changes.

A resource present in origin and deleted from upstream will be deleted from local.
A resource missing from origin and added in upstream will be added to local.
A resource only in local will be kept without changes.
A resource in both upstream and local will be merged into local.

On the field level, the rules differ based on the type of field.

For scalars and non-associative lists:

If the field is present in either upstream or local and the value is null, remove the field from local.
If the field is unchanged between upstream and local, leave the local value unchanged.
If the field has been changed in both upstream and local, update local with the value from upstream.

For mappings:

If the field is present in either upstream or local and the value is null, remove the field from local.
If the field is present only in local, leave the local value unchanged.
If the field is not present in local, add the delta between origin and upstream as the value in local.
If the field is present in both upstream and local, recursively merge the values between local, upstream and origin.

For associative lists:

If the field is present in either upstream or local and the value is null, remove the field from local.
If the field is present only in local, leave the local value unchanged.
If the field is not present in local, add the delta between origin and upstream as the value in local. ---------------
If the field is present in both upstream and local, recursively merge the values between local, upstream and origin.

# TODO

null values
input should not be modified, but memory may be shared
removing the value vs null
    what happens with deployment.spec.replica and the default value (thinking interaction with hpa)
