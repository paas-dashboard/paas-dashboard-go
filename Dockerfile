FROM shoothzj/compile:go AS build
COPY . /opt/compile
WORKDIR /opt/compile
RUN go build -o paas-dashboard .

FROM shoothzj/base

WORKDIR /opt/paas-dashboard

COPY --from=build /opt/compile/paas-dashboard /opt/paas-dashboard/paas-dashboard

RUN wget -q https://github.com/paas-dashboard/paas-dashboard-portal-angular/releases/download/latest/paas-dashboard-portal.tar.gz && \
    tar -xzf paas-dashboard-portal.tar.gz && \
    rm -rf paas-dashboard-portal.tar.gz

EXPOSE 11111

CMD ["/usr/bin/dumb-init", "/opt/paas-dashboard/paas-dashboard"]
