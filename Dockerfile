FROM scratch

# This image is meant to be very tiny and only execute purslane. No shell, no nothing. Therefore, the entrypoint below works just fine to be designed to only run purslane. No need to even add purslane to the PATH. 

COPY purslane /
ENTRYPOINT ["/purslane"] 