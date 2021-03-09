load("@io_bazel_rules_go//go:def.bzl", "go_test")

go_test(
  name = "go_default_testing",
  srcs = glob(['*.go']),
  tags = ['selftest'],
  deps = [
    "//staging/src/depops.io/api:go_default_library",
  ]
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
