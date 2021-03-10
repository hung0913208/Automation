load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
  name = "dataplan",
  srcs = ["dataplan.proto"],
  visibility = ["//visibility:public"],
)

go_proto_library(
  name = "",
  compiler = "@io_bazel_rules_go//proto:go_grpc",
  importpath = "devops.io/cloud/proto",
  proto = ":dataplan",
  visibility = ["//visibility:public"],
)

filegroup(
  name = "package-srcs",
  srcs = glob(["**"]),
  tags = ["automanaged"],
  visibility = ["//visibility:private"],
)

filegroup(
  name = "all-srcs",
  srcs = [
    ":package-srcs",
  ],
  tags = ["automanaged"],
  visibility = ["//visibility:public"],
)
